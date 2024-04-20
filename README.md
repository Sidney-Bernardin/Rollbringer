# STILL IN DEVELOPMENT

SELECT * FROM 
    (SELECT username FROM users) users
NATURAL FULL JOIN
    (SELECT name FROM pdfs WHERE owner_id = '7438abb6-6135-44de-a9aa-e755ce375505') pdfs
LIMIT 200
