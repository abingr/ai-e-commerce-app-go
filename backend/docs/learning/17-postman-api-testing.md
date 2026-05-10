# Postman API Testing

Postman is a GUI tool for sending HTTP requests and inspecting responses.

Phase 4 adds Postman files:

```text
postman/ai-e-commerce-app-go.postman_collection.json
postman/ai-e-commerce-local.postman_environment.json
```

Import both files into Postman.

Use the local environment when running requests. It stores variables such as:

- `base_url`
- `customer_email`
- `admin_email`
- `customer_token`
- `admin_token`
- `product_id`

The collection includes scripts that save tokens and product IDs automatically after successful responses.

Recommended flow:

1. Register Customer
2. Login Customer
3. Customer Cannot Create Product
4. Promote admin user in PostgreSQL
5. Login Admin
6. Create Product
7. Update Product
8. Delete Product

Backend engineer note:

- Postman is useful for exploratory testing.
- Automated tests still matter because Postman collections are not a replacement for CI.
- Keep example secrets local and avoid committing real tokens.

CV talking point:

> Added Postman collection and environment files to document and manually verify authentication and admin API workflows.
