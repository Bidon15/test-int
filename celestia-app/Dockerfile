# we need our own dockerhub org account
# below path is for testing the flow with @Bidon15 account
ARG branch=0.1.0
FROM celestiaorg/celestia-app:$branch 


USER root

RUN apk add --update bash curl jq
ENV NODE_NAME "user1"
# ports for p2p and rpc communication
EXPOSE 26656 26657 1317 9090

ADD capp.sh /capp.sh
RUN chmod +x /capp.sh
ENTRYPOINT ["/capp.sh"]