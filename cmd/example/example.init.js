window.addEventListener("load", function load(event){

    var results = document.getElementById("results");
    
    var draw_results = function(rsp){
	var data = JSON.parse(rsp);
	var str_data = JSON.stringify(data, "", " ");

	results.innerText = str_data;
    };
    
    // https://github.com/sfomuseum/go-http-wasm
    // https://github.com/sfomuseum/go-http-wasm/blob/main/static/javascript/sfomuseum.wasm.js
    
    sfomuseum.wasm.fetch("wasm/sfomuseum_airfield.wasm").then(rsp => {
	
	var airport_button = document.getElementById("airport_code_button");
	airport_button.removeAttribute("disabled");
	
	var airline_button = document.getElementById("airline_code_button");	
	airline_button.removeAttribute("disabled");

	var aircraft_button = document.getElementById("aircraft_code_button");	
	aircraft_button.removeAttribute("disabled");
	
	airport_button.onclick = function(){

	    var code_el = document.getElementById("airport_code");
	    var code = code_el.value;

	    sfomuseum_lookup_airport(code).then((rsp) => {
		draw_results(rsp);
	    }).catch((err) => {
		console.log(err);
	    });

	    
	    return false;
	};

	airline_button.onclick = function(){

	    results.innerHTML = "";
	    
	    var code_el = document.getElementById("airline_code");
	    var code = code_el.value;

	    sfomuseum_lookup_airline(code).then((rsp) => {
		draw_results(rsp);
	    }).catch((err) => {
		console.log(err);
	    });

	    return false;
	};

	aircraft_button.onclick = function(){

	    results.innerHTML = "";
	    
	    var code_el = document.getElementById("aircraft_code");
	    var code = code_el.value;

	    sfomuseum_lookup_aircraft(code).then((rsp) => {
		draw_results(rsp);
	    }).catch((err) => {
		console.log(err);
	    });

	    return false;
	};
	
    }).catch ((err) => {
	console.log("SAD", err);
    });
    
});
