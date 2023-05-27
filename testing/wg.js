// This is generated code from web.go, please be carefull editing this.
// Any changes made will be lost after a re-generate.
// Thank you for using web.go, more information here: https://github.com/Yadiiiig/web.go

const args = new Map();
const headers = new Headers();

function init() {
    headers.append("Content-Type", "application/json");

    fetchArgs()
        .then(() => {
            for (let [k, v] of args) {
                console.log(k, v);
                document.getElementById(k).getElementById = v;
            }
        })
        .catch(error => {
            console.error('Error occurred while fetching API data:', error);
        });
}

function render(data) {
    Object.entries(data).forEach(([responseKey, responseValue]) => {
        args.set(responseKey, responseValue);
    });

}

function fetchArgs(route) {
    return new Promise((resolve, reject) => {
        fetch("localhost:8080/"+route)
            .then(response => response.json())
            .then(data => {
                render(data)
                resolve(); // Resolve the promise after storing the data in the map
            })
            .catch(error => {
                console.error('Error occurred while sending API request:', error);
                reject(error); // Reject the promise if an error occurs
            });
    });
}


function remove() {
	let opts = {
        method: 'POST',
        headers: headers,
        body: [{'user.id': args['user.id']}, {' counter': args[' counter']}],
    };
        fetch("localhost:8080/foo/remove", opts)
            .then(response => response.text())
            .then(data => render(data))
        .catch(error => console.log('error', error));
    }
