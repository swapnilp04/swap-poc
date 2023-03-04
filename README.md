
## REGISTER
```
curl -v -H 'Content-Type: application/json' http://localhost:8080/register  -u "swapnil:mastahey" -d '{"username": "swap", "password": "admin123", "confirm_password": "admin123"}'
>> {"message":"user created successfully","user":{"id":1,"username":"swap"}}
```

## LOGIN
```
curl -v -H 'Content-Type: application/json' http://localhost:8080/login  -d '{"username": "swap", "password": "admin123"}'

{"message":"user loggedin successfully","token":"9848e3be-ba4b-11ed-b3a6-a660aea45daa"}
```

## UPDATE PASSWORD
```
curl -v -XPUT -H 'Content-Type: application/json' http://localhost:8080/updateUser  -H 'token: 9848e3be-ba4b-11ed-b3a6-a660aea45daa' -d '{"password": "admin123", "confirm_password": "admin123"}'
```

## LOGOUT 
```
curl -XDELETE -H 'Content-Type: application/json' http://localhost:8080/logout -H 'token: 0c145db4-ba5b-11ed-aadf-a660aea45daa'
>> {"message":"user logged out"}
```

## CREATE STUDENT

```
curl -XPOST -H 'Content-Type: application/json' http://localhost:8080/students -H 'token: 9160540c-ba5e-11ed-a389-a660aea45daa' -d '{"first_name": "Student One", "last_name": "Student Last", "age": 11, "content_number": 8888888888}'

>> {"message":"student created","student":{"id":2,"first_name":"Student Second","last_name":"Student Second","age":11,"phone_number":0}}
```

## UPDATE STUDENT
``` 
curl -XPUT -H 'Content-Type: application/json' http://localhost:8080/students/2 -H 'token: 9160540c-ba5e-11ed-a389-a660aea45daa' -d '{"first_name": "Student First2"}'
>> {"message":"student updated","student":{"id":1,"first_name":"Student First1","last_name":"Student Last","age":10,"phone_number":0}}
```


## DELETE STUDENT
```
curl  -XDELETE -H 'Content-Type: application/json' http://localhost:8080/students/2 -H 'token: 9160540c-ba5e-11ed-a389-a660aea45daa' 
>>> {"message":"student deleted successfully"}
```

## ALL
```
curl -XGET -H 'Content-Type: application/json' http://localhost:8080/students -H 'token: 9160540c-ba5e-11ed-a389-a660aea45daa'

>> [{"id":1,"first_name":"Student First","last_name":"Student Last","age":10,"phone_number":0},{"id":2,"first_name":"Student Second","last_name":"Student Second","age":11,"phone_number":0},{"id":3,"first_name":"Student One","last_name":"Student Last","age":11,"phone_number":0}]
```