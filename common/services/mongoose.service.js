const mongoose = require('mongoose');
let count = 0;

const connectWithRetry = async () => {
    console.log('MongoDB connection with retry')
    await mongoose.connect('mongodb://localhost/insurance-demo', {
    useNewUrlParser: true,
    useUnifiedTopology: true
    }).then(()=>{
        console.log('MongoDB is connected')
    }).catch(err=>{
        console.log('MongoDB connection unsuccessful, retry after 5 seconds. ', ++count);
        setTimeout(connectWithRetry, 5000)
    })
};

connectWithRetry();

exports.mongoose = mongoose;
