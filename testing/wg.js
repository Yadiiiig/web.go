
		// This is generated code from web.go, please be carefull editing this.\n		
		// Any changes made will be lost after a re-generate.\n
		// Thank you for using web.go, more information here: https://github.com/Yadiiiig/web.go\n
		
		const args = new Map();
		const headers = new Headers();
		let placeholders;

		function init(route) {
			headers.append("Content-Type", "application/json");
			placeholders = document.querySelectorAll("[data-token]");

			fetchArgs(route)
				.then(() => {
					render()
				})
				.catch(error => {
					console.error('Error occurred while fetching API data:', error);
				}
			);
		}

		function render() {
			placeholders.forEach(placeholder => {
				const token = placeholder.dataset.token;
				const replacementValue = args.get(token) || "";
				placeholder.innerHTML = replacementValue;
				console.log(token, args.get(token))
			});
		}

		function store(data) {
			Object.entries(data).forEach(([responseKey, responseValue]) => {
				args.set(responseKey, responseValue);
			});

		}

		function fetchArgs(route) {
			return new Promise((resolve, reject) => {
				fetch("http://127.0.0.1:8080/"+route)
					.then(response => response.json())
				.then(data => {
					store(data)
					resolve(); // Resolve the promise after storing the data in the map
				})
				.catch(error => {
					console.error('Error occurred while sending API request:', error);
					reject(error); // Reject the promise if an error occurs
				});
			});
		}

		function setParams(keys) {
			const selected = {};

			for (const key of keys) {
				if (args.has(key)) {
					selected[key] = args.get(key);
				}
			}

			return selected;
		}




		function remove(e) {
			e.preventDefault();	
			const selected = setParams(["user.id)"]);
	let opts = 
		{
			method: 'POST',
	        headers: headers,
		    body: JSON.stringify(selected),
		}
	;
        fetch("http://127.0.0.1:8080/foo/remove", opts)
            .then(response => response.json())
            .then(data => {
				console.log(data);
				store(data);
				render(data);
			})
        .catch(error => console.log('error', error));
    
		}
	
