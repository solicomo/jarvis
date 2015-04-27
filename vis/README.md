# Jarvis

J.A.R.V.I.S

## Startup

0. install tools

  0.1 install nodejs

      # wget http://nodejs.org/dist/v0.12.2/node-v0.12.2-linux-x64.tar.gz
      # tar -xvf node-v0.12.2-linux-x64.tar.gz
      # mv node-v0.12.2-linux-x64.tar.gz /usr/local/nodejs
      # ln -s /usr/local/nodejs/bin/* /usr/local/bin/

  0.2 install bower

      # npm install -g bower
      
1. install assets

    $ bower install

2. initate database

    $ sqlite3 app/data/vis.db
    > .read app/data/init.sql
