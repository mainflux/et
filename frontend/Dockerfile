FROM node:20.1.0

WORKDIR /usr/src/app/callhome-ui

COPY package*.json ./

RUN ls

RUN npm install -g @angular/cli @angular-devkit/build-angular && npm install

EXPOSE 4200

CMD ["ng", "serve", "--host", "0.0.0.0"]
