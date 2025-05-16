REVOKE ACCESS ON DATABASE `user-relation` FROM service_user;

REVOKE ALL GRAPH PRIVILEGES ON GRAPH `user-relation` FROM service_user;

REVOKE CREATE NEW NODE LABEL ON DATABASE `user-relation` FROM service_user;
REVOKE CREATE NEW RELATIONSHIP TYPE ON DATABASE `user-relation` FROM service_user;
REVOKE CREATE NEW PROPERTY NAME ON DATABASE `user-relation` FROM service_user;

DROP ROLE service_user;
