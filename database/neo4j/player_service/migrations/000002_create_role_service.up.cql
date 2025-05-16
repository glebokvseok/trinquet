CREATE ROLE service_user;

GRANT ACCESS ON DATABASE `user-relation` TO service_user;

GRANT ALL GRAPH PRIVILEGES ON GRAPH `user-relation` TO service_user;

GRANT CREATE NEW NODE LABEL ON DATABASE `user-relation` TO service_user;
GRANT CREATE NEW RELATIONSHIP TYPE ON DATABASE `user-relation` TO service_user;
GRANT CREATE NEW PROPERTY NAME ON DATABASE `user-relation` TO service_user;
