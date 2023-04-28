# Copyright wasilibs authors
# SPDX-License-Identifier: MIT

FROM ghcr.io/webassembly/wasi-sdk:wasi-sdk-20

RUN apt-get -y update && apt-get -y install curl

# Enables export of aligned_alloc
# TODO(anuraaga): It seems like this should be set automatically, consider debugging further.
ENV CFLAGS -O3 -D__USE_ISOC11 ${CFLAGS}

RUN mkdir -p /mimalloc && curl -L https://github.com/microsoft/mimalloc/archive/d1fff1119a52e15c1c3807efe1077024c39fe70e.tar.gz | tar -xz --strip-components 1 -C /mimalloc
WORKDIR /mimalloc
RUN mkdir -p out/release && cd out/release && cmake ../.. -DMI_BUILD_SHARED=off -DMI_BUILD_TESTS=off && make && ${RANLIB} libmimalloc.a

CMD ["cp", "./out/release/libmimalloc.a", "/out/libmimalloc.a"]
