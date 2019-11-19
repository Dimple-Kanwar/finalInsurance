// Get Weather by coodinates
let appId  = 'appid=b662848ce4c60f1334a86073c58f27de';
let url = 'http://api.openweathermap.org/data/2.5/weather?'
let lat="lat=12.98";
let lon="lon=77.6";
let units  = '&units=metric'; 
var request = require('request');
url = url+lat+"&"+lon+"&"+appId;
console.log("url: ", url)
request(url, function (error, response, body) {
    body = JSON.parse(body);
    if(error && response.statusCode != 200){
      throw error;
    }
  console.log("body: ", body);
});
     