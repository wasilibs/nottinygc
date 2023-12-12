# Copyright wasilibs authors
# SPDX-License-Identifier: MIT

FROM ghcr.io/webassembly/wasi-sdk:wasi-sdk-20

RUN apt-get -y update && apt-get -y install curl patch

RUN mkdir -p /bdwgc && curl -L https://github.com/ivmai/bdwgc/archive/a753bcb459649d20331d9c9f1dba94588ef56829.tar.gz | tar -xz --strip-components 1 -C /bdwgc
WORKDIR /bdwgc
ADD bdwgc.patch bdwgc.patch
RUN patch -p1 < bdwgc.patch
RUN ./autogen.sh

ENV CFLAGS -O3 ${CFLAGS}

# host is required by configure but not used so set it arbitrarily
RUN ./configure --disable-threads --disable-shared --disable-gcj-support --disable-java-finalization --disable-atomic-uncollectible --host=i686-pc-linux-gnu
RUN make

CMD ["sh", "-c", "cp /bdwgc/.libs/libgc.a /usr/lib/llvm-${LLVM_VERSION}/lib/clang/${LLVM_VERSION}/lib/wasi/libclang_rt.builtins-wasm32.a /out/"]
