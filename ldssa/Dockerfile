FROM ubuntu
RUN apt-get update && apt-get upgrade -y
#RUN apt-get update && apt-get install -y wget libc6:i386 libncurses5:i386 libstdc++6:i386 lib32z1 libbz2-1.0:i386
RUN apt-get update && apt-get install -y wget
RUN wget -q https://dl.google.com/android/repository/sdk-tools-linux-3859397.zip

CMD /bin/bash
