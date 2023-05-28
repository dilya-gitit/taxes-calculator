# taxes-calculator

Consists of two parts:
1.  Microservice to populate the database with taxation percentages for different years and for different social categories of workers
2.  A REST API Endpoint to calculate salary tax deductions for “ТОО на упрощенном режиме”. https://mybuh.kz/useful/calc_uproshenka/

#### Built with
- Go
- Gin
- gorm
- Postgresql

## Prerequisites

The following are needed to run this web application:

- [Docker](https://docs.docker.com/install/) (version 20.10)
- [Docker compose](https://docs.docker.com/compose/install/) (version 1.29)


## Install

Once the prerequisites are installed, execute the following commands from the project's root:
```bash
docker-compose up
```
This command will create containers for the Go server app and Postgres database.


## Usage

You can access the API endpoint by pasting the following URL in Postman: 
```bash
http://localhost:8080/post
```
It accepts following data as a raw JSON string:
```bash
salary(float,required)
year (int, required)
deduction 
is_staff_member(bool,required)
is_resident (bool, required)
social_statuses([]string,optional)
```
You can also choose social categories, to which your worker applies:
```bash
pensioner - "Пенсионер"
oppv - "Получатель ОППВ"
mother - "Многодетная мать"
disabled - "Инвалид"(1, 2, 3 группы)
disabled_perm - "Инвалид"(1, 2, 3 группы бессрочно)
student - "Студент"
asthub - "Сотрудник AstanaHub/МФЦА"
```
Output of the request is a JSON string with the amounts of all taxes applied and the net salary:
```bash
{
    "IPN": 3970,
    "OPV": 10000,
    "VOSMS": 2000,
    "net_salary": 84030
}
```
### Try running requests in Postman
[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/26583558-decc4630-3ff1-4e24-938d-8131a6453ee8?action=collection%2Ffork&collection-url=entityId%3D26583558-decc4630-3ff1-4e24-938d-8131a6453ee8%26entityType%3Dcollection%26workspaceId%3D14c6c675-d72c-4fc3-a26e-371fa7171831)

## Appendix
<img width="1013" alt="Screenshot 2023-05-28 at 08 12 11" src="https://github.com/dilya-gitit/taxes-calculator/assets/73358154/43ae14df-a166-46fe-a9c4-e628841bd7b0">






