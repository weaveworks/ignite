FROM node

RUN apt-get update && apt-get install -y git

RUN npm install gulp -g
RUN git clone https://github.com/github-tools/github-release-notes
WORKDIR /github-release-notes
RUN npm install
RUN gulp build && ln -s /github-release-notes/bin/gren.js /usr/local/bin/gren 
