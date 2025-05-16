exec 2>&1

scp -i ~/.ssh/yc \
    ~/GolandProjects/trinquet/.env \
    ~/GolandProjects/trinquet/tools/init_worker_node.sh \
    root@worker:/trinquet/
