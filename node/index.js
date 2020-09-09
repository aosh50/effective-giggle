var app = require('./app.js');


(async () => {

    const args = process.argv;
    // Default to file containing initial list
    var dataFile = "data/employee_list.csv";
    if (args.length > 2) {
        dataFile = args[2];
    }
    app.buildHierarchicalData(
        dataFile, 
        (trees, invalidEmployees) => {
            trees.forEach(tree => {
                var result = app.printTree(tree);
                console.table(result);
                invalidEmployees.forEach(e => console.log(`${e.name} has invalid manager id - ${e.manager_id}`));

            });
        },
        (error) => {
            console.error(error);
        });

    
})();