# iinservice

**REST API** for validating IINs and storing/retrieving people’s data,  
plus a **loader** for stress-testing with configurable concurrency.

---

## Tech Stack

- **Golang** (1.20)  
- **Chi** router & `go-chi/render`  
- **MongoDB** (6.0)  
- **Docker & Docker Compose**

---

## ⚙️ Configuration

### Environment Variables

| Variable           | Default                            | Description                                          |
|--------------------|------------------------------------|------------------------------------------------------|
| DATASTORE_URL      | `mongodb://localhost:27017`        | MongoDB connection URI                               |
| DATASTORE_DB       | `iinservice`                       | MongoDB database name                                |
| PORT               | `8080`                             | HTTP server port                                     |
| WORKERS            | `10`                               | Number of goroutines for loader                      |
| TOTAL_REQUESTS     | `100`                              | Total requests to send in loader                     |
| SERVICE_URL        | `http://app:8080`                  | Base URL for loader to hit (in Compose network)      |

---

## 🚀 Запуск с Docker Compose

Поднимаем MongoDB, API и loader одной командой:

```bash
make compose-up
```

Чтобы остановить и удалить контейнеры:

```bash
make compose-down
```

Если нужно запустить только стресс-тест:

```bash
make compose-loader
```

---

## 🛠️ Endpoints

### 1. Validate IIN  
```http
POST /iin_check
Content-Type: application/json

{ "iin": "770708389324" }
```
**Response**  
- `200 OK`  
  ```json
  { "valid": true, "date": "1997-03-05T00:00:00Z", "gender": "female" }
  ```
- `400 Bad Request`  
  ```json
  { "error": "IIN должен содержать ровно 12 цифр" }
  ```

---

### 2. Create Person  
```http
POST /people/info
Content-Type: application/json

{
  "iin":   "990422351053",
  "name":  "Иван Иванов",
  "phone": "+71234567890"
}
```
**Response**  
- `201 Created`  
  ```json
  { "iin": "990422351053" }
  ```
- `400 Bad Request` (already exists or invalid IIN/phone)

---

### 3. Get Person by IIN  
```http
GET /people/info/{iin}
```
**Response**  
- `200 OK`  
  ```json
  {
    "iin":   "990422351053",
    "name":  "Иван Иванов",
    "phone": "+71234567890"
  }
  ```
- `404 Not Found`  

---

### 4. Search by Name Part  
```http
GET /people/info/phone/{namePart}
```
**Response**  
- `200 OK` + JSON array of matches (or `[]` if none)  
  ```json
  [
    { "iin": "...", "name": "...", "phone": "..." },
    …
  ]
  ```

---

## 🔥 Stress Test Loader

Когда вы запускаете через Compose, `loader` автоматически пошлёт параллельные POST-запросы к вашему API:

```bash
make compose-loader
```



---

## 📖 License

MIT © Your Name
