# Stage 1
FROM node:16-alpine as builder 

WORKDIR /app 

COPY package*.json /app/ 

RUN npm install 

copy . /app

RUN npm run build 


# stage 2
FROM nginx:1.19.0 

WORKDIR /usr/share/nginx/html

# install vim 
RUN apt update && apt install vim -y

# replace default.conf in nginx conf files
COPY ./default.conf /etc/nginx/conf.d/

# remving default files in html dir
RUN rm -rf ./*

# copying the build/compiled files
COPY --from=builder /app/build .

# run nginx globally with deamon off to run nginx in foregroud
ENTRYPOINT ["nginx", "-g", "daemon off;"]

