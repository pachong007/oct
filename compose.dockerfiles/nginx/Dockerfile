FROM nginx

# set timezome
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

#RUN ln -s /usr/share/nginx/html/craw/resources/ /usr/share/nginx/html/one-cartoon/public/resources