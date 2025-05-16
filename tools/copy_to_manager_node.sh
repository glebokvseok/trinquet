exec 2>&1

ssh root@manager -i ~/.ssh/yc \
    "[[ -d /trinquet/conf/nginx/lib ]] || \
        sudo mkdir -p /trinquet/conf/nginx/lib && \
        sudo chmod 0777 /trinquet && \
        sudo chmod 0777 /trinquet/conf/nginx && \
        sudo chmod 0777 /trinquet/conf/nginx/lib"

ssh root@manager -i ~/.ssh/yc \
    "[[ -d /trinquet/conf/neo4j ]] || \
        sudo mkdir -p /trinquet/conf/neo4j && \
        sudo chmod 0777 /trinquet && \
        sudo chmod 0777 /trinquet/conf/neo4j"

scp -i ~/.ssh/yc \
    ~/GolandProjects/trinquet/.env \
    ~/GolandProjects/trinquet/docker-compose.init.yaml \
    ~/GolandProjects/trinquet/docker-compose.yaml \
    ~/GolandProjects/trinquet/tools/init_manager_node.sh \
    root@manager:/trinquet/

scp -i ~/.ssh/yc \
    ~/GolandProjects/trinquet/services/nginx/config/init.conf \
    ~/GolandProjects/trinquet/services/nginx/config/prod.conf \
    root@manager:/trinquet/conf/nginx

scp -i ~/.ssh/yc -r \
    ~/GolandProjects/trinquet/services/nginx/config/library/* \
    root@manager:/trinquet/conf/nginx/lib

scp -i ~/.ssh/yc -r \
    ~/GolandProjects/trinquet/services/neo4j/config/prod.conf \
    root@manager:/trinquet/conf/neo4j
