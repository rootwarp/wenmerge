export const getDifficalty = (target, cb) => {
    console.log('getDifficalty');
    const url = `https://api-wenmerge.dsrvlabs.dev/difficulty?target=${target}`;

    fetch(url)
        .then(resp => resp.json())
        .then(data => {
            cb(data);
        })
        .catch(err => {
            console.log(err);
        });
}
