// var request = require('request');
// var NodeGeocoder = require('node-geocoder');

// let url='https://maps.googleapis.com/maps/api/geocode/json?address=';
// let API_KEY='b1a922289ad34d00bda61f29445e2a07'
let address = "Lyngkein, Hima Mawphlang, Mawmyrsiang, Meghalaya 793121"
// url = url+address+"&"+API_KEY;
// console.log("url: ", url);


 
// var options = {
//   provider: 'opencage',
//   apiKey: API_KEY // for Mapquest, OpenCage, Google Premier
// };

// var geocoder = NodeGeocoder(options);
 
// // Using callback
// geocoder.geocode('29 champs elysée paris', function(err, res) {
//   console.log(err);
//   console.log(res);
// });

// // request(url, function (error, response, body) {
// //     body = JSON.parse(body);
// //     if(error && response.statusCode != 200){
// //       throw error;
// //     }
// //   console.log("body: ", body);
// // });

const opencage = require('opencage-api-client');

opencage.geocode({q: address}).then(data => {
  //console.log(JSON.stringify(data));
  if (data.status.code == 200) {
    if (data.results.length > 0) {
      var place = data.results[0];
      //console.log(place.formatted);
      console.log(place.geometry.lat);
      console.log(place.geometry.lng);
      coordinates = place.geometry.lat+","+place.geometry.lng
      console.log(coordinates);
      return coordinates;
      //console.log(place.annotations.timezone.name);
    }
  } else if (data.status.code == 402) {
    console.log('hit free-trial daily limit');
    console.log('become a customer: https://opencagedata.com/pricing'); 
  } else {
    // other possible response codes:
    // https://opencagedata.com/api#codes
    console.log('error', data.status.message);
    return data.status.message;
  }
}).catch(error => {
  console.log('error', error.message);
  return data.message;
});

// ... prints
// Theresienhöhe 11, 80339 Munich, Germany
// { lat: 48.1341651, lng: 11.5464794 }
// Europe/Berlin
