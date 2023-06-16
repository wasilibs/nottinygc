# Copyright wasilibs authors
# SPDX-License-Identifier: MIT

FROM ghcr.io/webassembly/wasi-sdk:wasi-sdk-20

RUN apt-get -y update && apt-get -y install curl

RUN mkdir -p /bdwgc && curl -L https://github.com/ivmai/bdwgc/archive/0a1667bee7f2c004abda6ecd173cd80bfef8412c.tar.gz | tar -xz --strip-components 1 -C /bdwgc
WORKDIR /bdwgc
RUN ./autogen.sh

# While signals aren't actually used the header file is included unconditionally. We can enable the header file
# without linking the actual signals library because the functions aren't called.
ENV CFLAGS -O3 -D_WASI_EMULATED_SIGNAL ${CFLAGS} -g

# host is required by configure but not used so set it arbitrarily
RUN ./configure --disable-threads --disable-shared --disable-gcj-support --disable-java-finalization --disable-atomic-uncollectible --host=i686-pc-linux-gnu
RUN make

CMD ["sh", "-c", "cp /bdwgc/.libs/libgc.a /usr/lib/llvm-${LLVM_VERSION}/lib/clang/${LLVM_VERSION}/lib/wasi/libclang_rt.builtins-wasm32.a /out/"]
