<div align="center">
      <h1>Basic Trade API</h1>
</div>

### Author
Asrie Annisa Pratiwi

### Description
A simple API Product & Variant with CRUD function and authentication (login & register)

### Table
| Table Name  | Column |
| -----------  | ----------- | 
|Admins  |id, uuid, name, email, password, created_at, updated_at |
|Products |id, uuid, name, image_url, admin_id, created_at, updated_at |
|Variants |id, uuid, variant_name, quantity, product_id, created_at, updated_at |

### Features
This API developed with Go version go1.20.6

### Tech Used
Golang, MySQL, Postman

### Getting Start:
Before you running the program, make sure you've set credential database, cloudinary, and port .env
| Environment | 
| ----------- | 
|HOST= |
|DB_USER= |
|DB_PASSWORD= |
|DB_NAME= |
|DB_PORT= |
|CLOUDINARY_CLOUD_NAME= |
|CLOUDINARY_API_KEY= |
|CLOUDINARY_API_SECRET= |
|CLOUDINARY_UPLOAD_FOLDER= |
|PORT= | 

### Run the program
https://basictradeasrie-production.up.railway.app/

<br/>you can read documentation collection Basic Trade API on link below
<br/><a href="https://documenter.getpostman.com/view/3449886/2s9YC7SBkc" target="_blank">https://documenter.getpostman.com/view/3449886/2s9YC7SBkc</a>

or if you will run on localhost you just used <b>go run .</b> on terminal vs code

### API Route List
| Method | URL | Description |
| ----------- | ----------- | ----------- | 
| POST | https://basictradeasrie-production.up.railway.app/auth/register  | Register User |
| POST | https://basictradeasrie-production.up.railway.app/auth/login  | Login User |
| GET | https://basictradeasrie-production.up.railway.app/products  | Get All Products |
| GET | https://basictradeasrie-production.up.railway.app/products/:uuid  | Get Product by UUID |
| POST | https://basictradeasrie-production.up.railway.app/products  | Create Product |
| PUT | https://basictradeasrie-production.up.railway.app/products/:uuid | Update Product |
| DELETE | https://basictradeasrie-production.up.railway.app/products/:uuid  | Delete Product |
| GET | https://basictradeasrie-production.up.railway.app/products/variants  | Get All Variants |
| GET | https://basictradeasrie-production.up.railway.app/products/variants/:uuid  | Get Variant by UUID |
| POST | https://basictradeasrie-production.up.railway.app/products/variants  | Create Variant (Add 'uuid product' params in request body for the create variant)|
| PUT | https://basictradeasrie-production.up.railway.app/products/variants/:uuid | Update Variant |
| DELETE | https://basictradeasrie-production.up.railway.app/products/variants/:uuid  | Delete Variant |


