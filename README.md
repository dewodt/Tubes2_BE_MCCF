# Wikirace Back-End

## Penjelasan Algoritma IDS dan BFS

Pada implementasi BFS kali ini Pages dianggap sebagai nodes, lalu jarak dari node awal ke suatu node disimpan dalam array jarak yang awal mula nya diset maxInt selain startURL yang diset 0 . Graf direpresentasikan dalam bentuk adjacency list yang diimplementasikan dalam map yang key nya merupakan node dan valuenya adalah array of string. Total pengunjungan node disimpan dalam variabel traversed. Antrian node yang akan diproses disimpan dalam queue lalu solusi disimpan dalam solution tree dengan nama parent. Solution tree ini yang nantinya akan dikunjungi menggunakan DFS untuk menampilkan solusi.

Pada implementasi IDS, pages bisa dianggap sebagai nodes, dan link-link yang ada pada nodes adalah children dari page tersebut. Penelusuran secara IDS melakukan penelusuran secara bertahap dengan meningkatkan kedalaman pencarian pada setiap iterasinya. Selain itu pada IDS juga dilakukan pengecekan ketika link yang sudah ditelusuri pada iterasi sebelumnya, maka link tersebut tidak perlu di scrapping lagi, dan hanya melakukan load dari cache. Penelusuran dengan IDS memanfaatkan rekursi yang cara kerja rekursi sendiri menggunakan stack secara tidak langsung. Penelusuran rekursi ini akan berhenti ketika parameter IDS startURL = endURL, lalu path dari startURL sampai endURL akan di append ke resultPath.

## Requirement program dan instalasi tertentu bila ada

- Docker 26.1.0
- Go version 1.21.6

## How To Run

1. Clone this repository

```bash
git clone https://github.com/dewodt/Tubes2_BE_MCCF
```

2. Go to `/src` directory

```bash
cd src
```

3. Make sure Docker Desktop is running

4. For development (with hot reload), run

- Build docker image (once)

```bash
docker compose -f "docker-compose.development.yml" build
```

- Run docker container

```bash
docker compose -f "docker-compose.development.yml" up
```

3. For production, run

- Build docker image (once)

```bash
docker compose -f "docker-compose.production.yml" build
```

- Run docker container

```bash
docker compose -f "docker-compose.production.yml" up
```

4. Visit [http://localhost:8080](http://localhost:8080)
5. To stop the process, find the container id

```bash
docker ps
```

6. Stop the process

```bash
docker stop <container_id>
```

## Author (Identitas Pembuat)

|   NIM    |         Nama         |
| :------: | :------------------: |
| 13522011 | Dewantoro Triatmojo  |
| 13522079 | Emery Fathan Zwageri |
| 13522095 | Rayhan Fadhlan Azka  |
