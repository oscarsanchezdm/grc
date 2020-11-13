FROM scratch

LABEL authors="Òscar Sánchez"

ADD bin/grc grc

EXPOSE 8000

CMD ["/grc"]
