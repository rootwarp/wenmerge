export const getDifficalty = (rpc, target, cb) => {
    console.log('getDifficalty');
    const url = `${rpc}/difficulty?target=${target}`;

    fetch(url)
        .then(resp => resp.json())
        .then(data => {
            cb(data);
        })
        .catch(err => {
            console.log(err);
        });
}
