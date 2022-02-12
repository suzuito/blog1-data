
```bash
# Upload modified articles
gsutil rsync -r articles gs://suzuito-minilla-blog1-article
gsutil rsync -r articles gs://suzuito-godzilla-blog1-article

# Delete articles
 gsutil rm gs://suzuito-minilla-blog1-article/2021-01-01-test.md
 gsutil rm gs://suzuito-godzilla-blog1-article/2021-01-01-test.md
```