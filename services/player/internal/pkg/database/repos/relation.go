package repos

import (
	"context"
	"github.com/google/uuid"
	"github.com/move-mates/trinquet/library/database/pkg/neo4j"
	apierrors "github.com/move-mates/trinquet/services/player/internal/pkg/errors"
	"github.com/move-mates/trinquet/services/player/internal/pkg/models"
	"github.com/pkg/errors"
	"time"
)

const (
	mergeUserNodeQuery = `MERGE (u:User {id: $id})`

	followUserQuery = `
		OPTIONAL MATCH (self:User {id: $self_id})
		OPTIONAL MATCH (user:User {id: $user_id})
		WITH self, user, 
			self IS NOT NULL AS self_exists,
			user IS NOT NULL AS user_exists,
			EXISTS { MATCH (self)-[:FOLLOW]->(user) } AS followership_existed
		MERGE (self)-[r:FOLLOW]->(user)
		ON CREATE SET 
			r.followed_on = $followed_on, 
			r.created_on = datetime().epochMillis
		RETURN self_exists, user_exists, followership_existed
	`

	unfollowUserQuery = `
		OPTIONAL MATCH (self:User {id: $self_id})
		OPTIONAL MATCH (user:User {id: $user_id})
		OPTIONAL MATCH (self)-[r:FOLLOW]->(user)
		DELETE r
		RETURN
		  	self IS NOT NULL AS self_exists,
		  	user IS NOT NULL AS user_exists, 
		  	r IS NOT NULL AS followership_existed
	`

	getFollowInfoQuery = `
		MATCH (self:User {id: $self_id})
		MATCH (user:User {id: $user_id})

		WITH self, user

		OPTIONAL MATCH (user)<-[:FOLLOW]-(follower)
		WITH self, user, collect(DISTINCT follower) AS followers
		
		OPTIONAL MATCH (user)-[:FOLLOW]->(following)
		WITH self, user, followers, collect(DISTINCT following) AS following
		
		OPTIONAL MATCH (self)-[s2u:FOLLOW]->(user)
		OPTIONAL MATCH (user)-[u2s:FOLLOW]->(self)

		RETURN 
  			size(followers) AS followers_count,
  			size(following) AS following_count,
			s2u IS NOT NULL AS following,
			u2s IS NOT NULL AS following_back
	`

	getFollowsQuery = `
		UNWIND $user_ids AS user_id

		MATCH (self:User {id: $self_id})
		MATCH (user:User {id: user_id})

		WITH self, user
		
		OPTIONAL MATCH (self)-[s2u:FOLLOW]->(user)
		OPTIONAL MATCH (user)-[u2s:FOLLOW]->(self)

		RETURN
			user.id AS user_id,
			s2u IS NOT NULL AS following,
			u2s IS NOT NULL AS following_back
	`

	getAllUserFollowerIDsQuery = `
		MATCH (user:User {id: $id})
		OPTIONAL MATCH (user)<-[:FOLLOW]-(follower:User)
		RETURN follower.id AS follower_id
	`

	getUserFollowersQuery = `
		OPTIONAL MATCH (user:User {id: $id}) 
		OPTIONAL MATCH (user)<-[r:FOLLOW]-(follower:User)
		WHERE r.followed_on < $cursor
		OPTIONAL MATCH (user)-[u2f:FOLLOW]->(follower)
		RETURN 
			   user IS NOT NULL as user_exists,
			   follower.id as follow_id,
			   r.followed_on as followed_on,
			   u2f IS NOT NULL AS is_mutual
		ORDER BY r.followed_on DESC
		LIMIT $limit
	`

	getUserFollowingsQuery = `
		OPTIONAL MATCH (user:User {id: $id}) 
		OPTIONAL MATCH (user)-[r:FOLLOW]->(following:User)
		WHERE r.followed_on < $cursor
		OPTIONAL MATCH (following)-[f2u:FOLLOW]->(user)
		RETURN 
			   user IS NOT NULL as user_exists,
			   following.id as follow_id,
			   r.followed_on as followed_on,
			   f2u IS NOT NULL AS is_mutual
		ORDER BY r.followed_on DESC
		LIMIT $limit
	`
)

type followDirection bool

const (
	follower  followDirection = true
	following followDirection = false
)

type RelationRepository interface {
	CreateUserNodeIfNotExists(ctx context.Context, id uuid.UUID) error
	FollowUser(ctx context.Context, selfID uuid.UUID, userID uuid.UUID, followedOn time.Time) error
	UnfollowUser(ctx context.Context, selfID uuid.UUID, userID uuid.UUID) error
	GetFollowInfo(ctx context.Context, selfID uuid.UUID, userID uuid.UUID) (models.FollowInfo, error)
	GetFollows(ctx context.Context, selfID uuid.UUID, userIDs []uuid.UUID) ([]models.Follow, error)
	GetAllFollowerIDs(ctx context.Context, userID uuid.UUID) ([]string, error)
	GetFollowers(ctx context.Context, userID uuid.UUID, sort models.FollowSort, followersCount int) (models.FollowSortResult, error)
	GetFollowing(ctx context.Context, userID uuid.UUID, sort models.FollowSort, followingCount int) (models.FollowSortResult, error)
}

type relationRepository struct {
	driver neo4j.Driver
	config neo4j.SessionConfig
}

func ProvideRelationRepository(
	driver neo4j.Driver,
	config neo4j.SessionConfig,
) RelationRepository {
	return &relationRepository{
		driver: driver,
		config: config,
	}
}

func (repo *relationRepository) CreateUserNodeIfNotExists(ctx context.Context, id uuid.UUID) (returnErr error) {
	session, ctx, cancel := neo4j.NewWriteSession(ctx, repo.driver, repo.config)
	defer cancel()
	defer neo4j.SessionFinalizer(ctx, session, &returnErr)

	_, err := session.Run(ctx, mergeUserNodeQuery, map[string]any{"id": id.String()})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *relationRepository) FollowUser(ctx context.Context, selfID uuid.UUID, userID uuid.UUID, followedOn time.Time) error {
	if selfID == userID {
		return apierrors.NewSelfFollowError()
	}

	params := map[string]any{
		"self_id":     selfID.String(),
		"user_id":     userID.String(),
		"followed_on": followedOn.UnixMilli(),
	}

	followerShipExisted, err := repo.makeFollowershipChange(ctx, followUserQuery, params)
	if err != nil {
		return errors.WithStack(err)
	}

	if followerShipExisted {
		return apierrors.NewFollowershipAlreadyExistsError(userID)
	}

	return nil
}

func (repo *relationRepository) UnfollowUser(ctx context.Context, selfID uuid.UUID, userID uuid.UUID) error {
	params := map[string]any{
		"self_id": selfID.String(),
		"user_id": userID.String(),
	}

	followershipExisted, err := repo.makeFollowershipChange(ctx, unfollowUserQuery, params)
	if err != nil {
		return errors.WithStack(err)
	}

	if !followershipExisted {
		return apierrors.NewFollowershipDoesNotExistError(userID)
	}

	return nil
}

func (repo *relationRepository) makeFollowershipChange(
	ctx context.Context,
	query string,
	params map[string]any,
) (
	returnFollowershipExisted bool,
	returnErr error,
) {
	session, ctx, cancel := neo4j.NewWriteSession(ctx, repo.driver, repo.config)
	defer cancel()
	defer neo4j.SessionFinalizer(ctx, session, &returnErr)

	res, err := session.Run(ctx, query, params)
	if err != nil {
		return false, errors.WithStack(err)
	}

	if !res.Next(ctx) {
		if res.Err() != nil {
			return false, errors.WithStack(res.Err())
		}

		return false, errors.Errorf("make followership change request returned no records")
	}

	values, err := getValues(res, []string{"self_exists", "user_exists", "followership_existed"})
	if err != nil {
		return false, err
	}

	if !values["self_exists"].(bool) {
		return false, apierrors.NewUserNodeNotFoundError(params["self_id"].(string))
	}

	if !values["user_exists"].(bool) {
		return false, apierrors.NewUserNodeNotFoundError(params["user_id"].(string))
	}

	return values["followership_existed"].(bool), nil
}

func (repo *relationRepository) GetFollowInfo(ctx context.Context, selfID uuid.UUID, userID uuid.UUID) (
	followInfo models.FollowInfo,
	returnErr error,
) {
	session, ctx, cancel := neo4j.NewReadSession(ctx, repo.driver, repo.config)
	defer cancel()
	defer neo4j.SessionFinalizer(ctx, session, &returnErr)

	res, err := session.Run(ctx, getFollowInfoQuery,
		map[string]any{
			"self_id": selfID.String(),
			"user_id": userID.String(),
		},
	)

	if err != nil {
		return followInfo, errors.WithStack(err)
	}

	if !res.Next(ctx) {
		if res.Err() != nil {
			return followInfo, errors.WithStack(res.Err())
		}

		return followInfo, errors.Errorf("count followers and following request returned no records")
	}

	values, err := getValues(res, []string{"followers_count", "following_count", "following", "following_back"})
	if err != nil {
		return followInfo, err
	}

	return models.FollowInfo{
		FollowersCount: values["followers_count"].(int64),
		FollowingCount: values["following_count"].(int64),
		Following:      values["following"].(bool),
		FollowingBack:  values["following_back"].(bool),
	}, nil
}

func (repo *relationRepository) GetFollows(ctx context.Context, selfID uuid.UUID, userIDs []uuid.UUID) (
	follows []models.Follow,
	returnErr error,
) {
	session, ctx, cancel := neo4j.NewReadSession(ctx, repo.driver, repo.config)
	defer cancel()
	defer neo4j.SessionFinalizer(ctx, session, &returnErr)

	var ids = make([]string, len(userIDs))
	for i := range userIDs {
		ids[i] = userIDs[i].String()
	}

	res, err := session.Run(ctx, getFollowsQuery,
		map[string]any{
			"self_id":  selfID.String(),
			"user_ids": ids,
		},
	)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	followsMap := make(map[uuid.UUID]models.Follow, len(userIDs))
	for res.Next(ctx) {
		if res.Err() != nil {
			return nil, errors.WithStack(res.Err())
		}

		values, err := getValues(res, []string{"user_id", "following", "following_back"})
		if err != nil {
			return nil, err
		}

		userID, err := uuid.Parse(values["user_id"].(string))
		if err != nil {
			return nil, errors.WithStack(err)
		}

		followsMap[userID] = models.Follow{
			ID:            userID,
			Following:     values["following"].(bool),
			FollowingBack: values["following_back"].(bool),
		}
	}

	follows = make([]models.Follow, len(userIDs))
	for i := range follows {
		follows[i] = followsMap[userIDs[i]]
	}

	return follows, nil
}

func (repo *relationRepository) GetAllFollowerIDs(ctx context.Context, userID uuid.UUID) (followerIDs []string, returnErr error) {
	session, ctx, cancel := neo4j.NewReadSession(ctx, repo.driver, repo.config)
	defer cancel()
	defer neo4j.SessionFinalizer(ctx, session, &returnErr)

	res, err := session.Run(ctx, getAllUserFollowerIDsQuery, map[string]any{"id": userID.String()})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if res.Err() != nil {
		return nil, errors.WithStack(res.Err())
	}

	for res.Next(ctx) {
		if res.Err() != nil {
			return nil, errors.WithStack(res.Err())
		}

		followerID, ok := res.Record().Get("follower_id")
		if !ok {
			return nil, errors.Errorf("failed to get follower_id field")
		}

		followerIDs = append(followerIDs, followerID.(string))
	}

	return followerIDs, nil
}

func (repo *relationRepository) GetFollowers(
	ctx context.Context,
	userID uuid.UUID,
	sort models.FollowSort,
	followersCount int,
) (models.FollowSortResult, error) {
	return repo.getFollows(ctx, userID, sort, followersCount, follower)
}

func (repo *relationRepository) GetFollowing(
	ctx context.Context,
	userID uuid.UUID,
	sort models.FollowSort,
	followingCount int,
) (models.FollowSortResult, error) {
	return repo.getFollows(ctx, userID, sort, followingCount, following)
}

func (repo *relationRepository) getFollows(
	ctx context.Context,
	userID uuid.UUID,
	sort models.FollowSort,
	followCount int,
	direction followDirection,
) (sortResult models.FollowSortResult, returnErr error) {
	if sort.Type != models.DateTimeDesc {
		return models.FollowSortResult{}, apierrors.NewUnsupportedFollowSortType(sort.Type)
	}

	session, ctx, cancel := neo4j.NewReadSession(ctx, repo.driver, repo.config)
	defer cancel()
	defer neo4j.SessionFinalizer(ctx, session, &returnErr)

	query := getUserFollowingsQuery
	if direction == follower {
		query = getUserFollowersQuery
	}

	res, err := session.Run(ctx, query,
		map[string]any{
			"id":     userID.String(),
			"cursor": sort.Cursor,
			"limit":  followCount,
		},
	)

	if err != nil {
		return models.FollowSortResult{}, errors.WithStack(err)
	}

	if res.Err() != nil {
		return models.FollowSortResult{}, errors.WithStack(res.Err())
	}

	newCursor := sort.Cursor
	var follows []models.Follow
	for res.Next(ctx) {
		if res.Err() != nil {
			return models.FollowSortResult{}, errors.WithStack(res.Err())
		}

		values, err := getValues(res, []string{"user_exists", "follow_id", "followed_on", "is_mutual"})
		if err != nil {
			return models.FollowSortResult{}, err
		}

		if !values["user_exists"].(bool) {
			return models.FollowSortResult{}, apierrors.NewUserNodeNotFoundError(userID.String())
		}

		followID, err := uuid.Parse(values["follow_id"].(string))
		if err != nil {
			return models.FollowSortResult{}, errors.WithStack(err)
		}

		newCursor = min(newCursor, values["followed_on"].(int64))

		isMutual := values["is_mutual"].(bool)
		follows = append(follows, models.Follow{
			ID:            followID,
			Following:     direction == following || isMutual && direction == follower,
			FollowingBack: direction == follower || isMutual && direction == following,
		})
	}

	return models.FollowSortResult{
		Follows: follows,
		Cursor:  newCursor,
	}, nil
}

func getValues(result neo4j.Result, fields []string) (map[string]any, error) {
	values := make(map[string]any)
	for _, field := range fields {
		value, ok := result.Record().Get(field)
		if !ok {
			return nil, errors.Errorf("failed to get %s field", field)
		}

		values[field] = value
	}

	return values, nil
}
