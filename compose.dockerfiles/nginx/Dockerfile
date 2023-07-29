FROM nginx

# set timezome
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN ln -s /usr/share/nginx/html/oct/www/crawler/resources/ /usr/share/nginx/html/oct/www/one-cartoon/public/resources