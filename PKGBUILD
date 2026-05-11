# Maintainer: Asep5K <asepdev.git@gmail.com>
pkgname=animein-git
pkgver=r6.e8dbf80
pkgrel=1
pkgdesc='CLI for watching anime from https://animeinweb.com.'
url=https://codeberg.org/Asep5K/animein
arch=('any')
license=('GPL3')
makedepends=('go' 'git')
provides=('animein')
conflicts=('animein')
source=("${pkgname}-git::git+${url}.git")
sha256sums=('SKIP')

pkgver() {
	cd "${srcdir}/${pkgname}"
	printf "r%s.%s" "$(git rev-list --count HEAD)" "$(git rev-parse --short HEAD)"
}

prepare() {
	export GOPATH="${srcdir}/gopath"
	go clean -modcache
}

build() {
	cd "${srcdir}/${pkgname}"
	export GOPATH="${srcdir}/gopath"
	export CGO_CPPFLAGS="${CPPFLAGS}"
	export CGO_CFLAGS="${CFLAGS}"
	export CGO_CXXFLAGS="${CXXFLAGS}"
	export CGO_LDFLAGS="${LDFLAGS}"
	export GOFLAGS="-buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw"

	go build -o build/ -v .
	go clean -modcache
}

package(){
	cd "${srcdir}/${pkgname}"
	install -Dvm755 "build/${pkgname}" -t "${pkgdir}/usr/bin"
	install -Dvm644 'LICENSE' "${pkgdir}/usr/share/licenses/${pkgname}/LICENSE"
}
