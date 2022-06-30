export const getDifficalty = (target) => {
    console.log('getDifficalty');
    const url = 'http://localhost:9090/difficulty';

    fetch(url)
    .then(result => {
        console.log(result);
    })
    .catch(err => {
        console.log(err);
    });

    return;
}