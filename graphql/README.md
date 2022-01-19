# Graphql

```graphql
query {
  gebouw(identificatie: "0303100000586320") {
    domein
    identificatie
    oorspronkelijkBouwjaar
    status
    type
    geregistreerdMet {
      beginGeldigheid
      eindGeldigheid
      afgeleidVan {
        versie
        beginGeldigheid
        eindGeldigheid
        beschrijft {
          domein
          identificatie
        }
      }
    }
  }
}
```

```json
{
    "data": {
        "gebouw": {
            "domein": "sor/gebouw",
            "identificatie": "0303100000586320",
            "oorspronkelijkBouwjaar": "1976",
            "status": "Bestaand",
            "type": "Pand",
            "geregistreerdMet": {
                "beginGeldigheid": "2012-11-09",
                "eindGeldigheid": null,
                "afgeleidVan": [
                    {
                        "versie": 2,
                        "beginGeldigheid": "2012-11-09",
                        "eindGeldigheid": null,
                        "beschrijft": {
                            "domein": "bag/pand",
                            "identificatie": "0303100000586320"
                        }
                    },
                    {
                        "versie": null,
                        "beginGeldigheid": "2018-08-09",
                        "eindGeldigheid": null,
                        "beschrijft": {
                            "domein": "bgt/pand",
                            "identificatie": "G0303.0979f330198c319ae05332a1e90a5e0b"
                        }
                    }
                ]
            }
        }
    }
}
```

```graphql
{
  gebouwCollectie(filter: { geometrie: { intersects: { fromWKT: "POLYGON ((175055.291 505776.506, 175060.893 505787.679, 175055.18 505790.543, 175056.524 505793.263, 175053.341 505794.851, 175049.36 505786.979, 175052.482 505785.4, 175049.928 505780.321, 175052.683 505778.936, 175052.234 505778.042, 175055.291 505776.506))" }}}) {
    nodes {
      domein
      identificatie
      status
      oorspronkelijkBouwjaar
      geometrie (srid: 9067 ){
        asWKB
      }
      geregistreerdMet{
        beginGeldigheid
      }
    }
  }
}
```

```json
{
  "data": {
    "gebouwCollectie": {
      "nodes": [
        {
          "domein": "sor/gebouw",
          "identificatie": "0303100000523550",
          "status": "Bestaand",
          "oorspronkelijkBouwjaar": "2011",
          "geometrie": {
            "asWKB": "ACAAAAMAACNrAAAAAQAAAAtAFrszdV6ND0BKRQaMbgRiQBa7SUfYwuBASkUJ1QbpOUAWuzNAWnRRQEpFCq6z9tdAFrs4fMkpVEBKRQt7WfmtQBa7LDa/r65ASkUL9Av4U0AWuxy1DuJoQEpFCaPPn91AFrsovtTWSkBKRQkrzbkxQBa7Hsv7nyZASkUHra+MW0AWuylrgKP3QEpFB0Rs7GdAFrsnq8CNpEBKRQcBKmZKQBa7M3VejQ9ASkUGjG4EYg=="
          },
          "geregistreerdMet": {
            "beginGeldigheid": "2012-11-09"
          }
        },
        {
          "domein": "sor/gebouw",
          "identificatie": "0303100000523553",
          "status": "Bestaand",
          "oorspronkelijkBouwjaar": "2011",
          "geometrie": {
            "asWKB": "ACAAAAMAACNrAAAAAQAAAAxAFrszdV6ND0BKRQaMbgRiQBa7Pz4MB9RASkUGF5wo60AWu0D+vKLNQEpFBlrLWztAFrtLnkGnnkBKRQXxiLtHQBa7VJrQwBBASkUHSsi3UEAWu1hJsemLQEpFB2DhMCBAFrth08laFUBKRQcCB5cjQBa7cVSa0MBASkUJU9nPckAWu2TtFJXkQEpFCc/zrvVAFrtfQXLVz0BKRQj7zdJtQBa7SUfYwuBASkUJ1QbpOUAWuzN1Xo0PQEpFBoxuBGI="
          },
          "geregistreerdMet": {
            "beginGeldigheid": "2012-11-09"
          }
        }
      ]
    }
  }
}
```
