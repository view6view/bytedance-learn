select `name`, `age`, `employee_number`
FROM `users`
where role <> "manager"
  AND age > 35
ORDER BY age DESC LIMIT 10
OFFSET 10 FOR
UPDATE