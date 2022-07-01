export const getDifficalty = (target, cb) => {
    console.log('getDifficalty');
    const url = `http://localhost:9090/difficulty?target=${target}`;

    fetch(url)
        .then(resp => resp.json())
        .then(data => {
            cb(data);
        })
        .catch(err => {
            console.log(err);
        });
}