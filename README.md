### Requirements
設計特定航空的訂機票的功能，規格：

1. 可以按照起點、目的地、日期查詢班機狀態

    ```bash
    curl localhost:9090/api/flights?departure_time=2022-01-01T08:00:00Z&seats=1&destination=Los Angeles&departure=New York&page=2&per_page=10
    ```

2. 採取分頁方式返回可用班機清單、價格以及剩餘機位

    ```bash
    curl localhost:9090/api/flights?departure_time=2022-01-01T08:00:00Z&seats=1&destination=Los Angeles&departure=New York&page=2&per_page=10
    ```

3. 航空公司有超賣的慣例，功能也需要考量到超賣的情境。

    When a user books or modifies a reservation, **optimistic locking** is used to check the remaining seats on the flight. This means that in high-concurrency scenarios, overbooking is allowed.

4. 設計表結構和索引、編寫主要程式碼、考慮大流量高並發情況 (可以使用虛擬碼實現)。

    The reason for using optimistic locking is to accommodate scenarios with **high read operations and low write operations**. Compared to pessimistic locking, optimistic locking helps prevent database contention and excessive waiting during high-traffic, high-concurrency situations.

### How To Run This Project

> Make Sure you have run the flight.sql in your mysql

Since the project is already use Go Module, I recommend to put the source code in any folder but GOPATH.

#### Run the Testing

```bash
$ make tests
```

#### Run the Applications

Here is the steps to run it with `docker-compose`

```bash

# Run the application
$ make up

# The hot reload will running

# Execute the call in another terminal
$ curl localhost:9090/api/flights
```
