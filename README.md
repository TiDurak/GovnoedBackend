<h1 align="center">
    <img src="https://user-images.githubusercontent.com/82606298/170775456-475ffa71-9cf9-4584-9723-b3917ae0aecc.svg" alt="DebilBot" border="0" height="30px"> 
    DebilBot
</h1>


# GovnoedPromoAPI

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)  

![Discord](https://img.shields.io/discord/439778332121497612?style=flat-square&logo=discord&logoColor=white&label=Discord%20Server)
![Website](https://img.shields.io/website?url=https%3A%2F%2Fgovnoed.de%2F&up_message=Visit&down_message=Not%20working&style=flat-square&label=govnoed.de)
![GitHub](https://img.shields.io/website?url=https%3A%2F%2Fgithub.com%2FTiDurak%2Fdebilbot&up_message=tidurak%2Fdebilbot&up_color=black&down_message=removed&style=flat-square&logo=github&label=GitHub)


## Description

This project is a refactored backend for the [govnoed.de](https://govnoed.de/) website.

The API uses an external SQLite database, communicates with the website through HTTP API endpoints, and generates promo codes for [debilbot](https://github.com/tidurak/debilbot).


## Project Status

The project is in its final stage. Future bug fixes and other minor fixes may still be made.


## Project Structure
```
├── cmd/
│   └── api/
│       └── main.go # API entry point
└── internal/
    ├── config/
    │   └── config.go # Configuration loading and validation
    ├── database/
    │   ├── database.go # SQLite connection
    │   ├── migrations/
    │   │   └── 001_create_redeem_keys.sql # Creates the redeem keys table
    │   │                             
    │   └── migrations.go # Database migration management
    ├── handler/
    │   ├── generateRequest.go # Key generation request structure
    │   ├── generateResponse.go # Key generation response structure
    │   ├── health.go # HTTP health check
    │   ├── promo.go # Promo key HTTP handlers
    │   └── response.go # HTTP response helpers
    ├── repository/
    │   └── promo.go # Promo key database operations
    └── service/
        └── promo.go # Promo key business logic
```

```mermaid
---
title: GovnoedPromoAPI
---
flowchart
	subgraph s1["GovnoedPromoAPI"]
		subgraph s6["Handler"]
			n10["127.0.0.1:8080/api/generate"]
		end
		subgraph s4["Service"]
			subgraph s5["Repository"]
				subgraph s7["Database"]
					n6
				end
				n7@{ shape: "diam", label: "if true" }
				n4@{ shape: "diam", label: "if false" }
				n1["now - created_at &lt; 12h"]
			end
			n9["key_hash"]
			n8["generate key"]
		end		
	end
	subgraph s3["DebilBot"]
		n5@{ shape: "lean-r", label: "is_used - edit to True / check" }
		
		n3["promo_keys.py"]
	end
	n2@{ shape: "cyl", label: "sqlite" }
	n2 ---|"Get values from db"| s3
	subgraph s2["Website + PHP"]
		
	end
	n5@{ shape: "lean-r", label: "check if is_used" }
	s3 ---|"is_used - edit to True"| n2
	style s1 fill:#CB6CE6,stroke:#8C52FF,stroke-width:2px
	n8 --- n9
	n1 --- n4
	n1 --- n7
	s2 ---|"discord_id"| s6
	n7 ---|"return error"| s6
	s6 ---|"return key, reward or error"| s2
	s6 --- s5
	n4 --- s6
	n9 --- s5
	n6["open database"]
	n4 --- s7
	s7 ---|"INSERT discord_id, created_at, key_hash"| n2
	s7 --- n1
	style s4 fill:#0CC0DF
	style s5 color:#00BF63,fill:#7ED957
	style s7 fill:#FFDE59
	style s6 fill:#FF5757
	style n10 fill:#FF66C4
	style s2 color:#FFFFFF,fill:#545454,stroke-width:0px
	style s3 color:#FFFFFF,fill:#545454,stroke-width:0px
	style n5 stroke-width:2px,stroke-dasharray:5 5
```


## FAQ

Q: Why was the backend rewritten?  
A: As of 23 August 2026, the current [govnoed.de](https://govnoed.de/get_key) backend is written in PHP. PHP was the best choice for my website at the time. However, the project quickly became difficult to scale and turned into an unreadable mess.


Q: Why Go instead of Rust, Django, or Node.js?  
A: The goal was to build a lightweight API because my VPS has very limited resources. The choice ultimately came down to Rust and Go. I chose Go because it keeps the code much simpler with only a small loss in performance. I also had no prior experience with Rust.
