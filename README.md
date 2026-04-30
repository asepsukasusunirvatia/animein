# ANIMEIN-CLI
**CLI buat nonton anime di [animeinweb](https://animeinweb.com/) via terminal.
No bloat, pure Go (Bahasa Kudanil), sat-set wat-wet.**

**Inspired by [animeku-cli](https://github.com/lucasbuilds/animeku-cli), but rewritten in Pure Go for faster compilation, a minimal footprint, and a focus on core features.**

---

## Requirement
- **[mpv](https://github.com/mpv-player/mpv)**	<- Termux ga butuh ini
- **[go 1.26.2++](https://go.dev/dl/ "kudanil err != nil")**
- **make** <- opsional


## Install

- **Menggunakan make (rekomendasi)**

```bash
git clone https://codeberg.org/Asep5K/animein.git
cd animein
make install
```
- **Manual (malas banget)**

```bash
git clone https://codeberg.org/Asep5K/animein.git
cd animein
go build -o animein main.go
mv animein ~/.local/bin
```
---

## Usage
- **Langsung pake judul di cli**

```bash
animein "Keikenzumi na Kimi to, Keiken Zero na Ore ga, Otsukiai suru Hanashi"
```

- **atau pake mode interactive tinggal panggil aja binary animein nya**
```bash
animein
```
---

