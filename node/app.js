const csv = require('fast-csv');
const _ = require('ramda');



module.exports = {
    Employee: class {
        constructor(id, name, manager_id = '') {
            this.id = id;
            this.name = name;
            this.manager_id = manager_id;
        }
     
    },
    
    TreeNode: class {
        constructor(value) {
            this.value = value;
            this.children = [];
            this.parent = null;
          }
    },

    toModel: function (r) {
        var e = new this.Employee(r.id, r.name);
        if (r.manager_id) {
            e.manager_id = r.manager_id;
        }
        return e;
    },

    getDepth: function getDepth(node) {
        
        var depth = 0;
        if (node.children) {
            node.children.forEach(function (d) {
                var tmpDepth = getDepth(d)
                if (tmpDepth > depth) {
                    depth = tmpDepth
                }
            })
        }
        return 1 + depth
    },


    buildNode: function(current, parent, list) {
    
        var node = new this.TreeNode(current);

        if (parent !== null) { // ie not the top node
            node.parent = parent;
        }
        
        // Does this employee manage anyone else? 
        var children = _.filter(e => e.manager_id === current.id, list);

        var childNodes = [];
        _.map(child => {
                childNode = this.buildNode(child, node, list);
                childNodes.push(childNode);
            }
        , children);

        node.children = childNodes;

        return node;
    },

    loadEmployeesFromFile: function (file, callback) {
        var employees = [];
        var csvRows = [];
        csv.parseFile(file, { headers: true })
            .on('error', error => console.error(error))
            .on('data', row => {
                csvRows.push(row);
            })
            .on('end', () => {
                csvRows.forEach(r => employees.push(this.toModel(r)));
                callback(employees);
            });
    },

    buildHierarchicalData: function (file, callback, errorCallback) {
        // Load the employees
        this.loadEmployeesFromFile(file, (employees) => { 

            if (employees.length === 0) {
                errorCallback("No employees found!");
                return;
            }
    
            // Find the top node(s) of our tree
            var topLevelEmployees = _.filter(e => _.isEmpty(e.manager_id), employees);
    
            if (topLevelEmployees.length === 0) {
                errorCallback("No CEO found");
                return;
            }

            var employeeIds = _.pluck('id', employees);
            var employeesWithInvalidManager = _.filter(e => !_.isEmpty(e.manager_id) && !_.contains(e.manager_id, employeeIds), employees);
    
            var trees = [];
            _.map(emp => {
                tree = this.buildNode(emp, null, employees);
                trees.push(tree)
            }, topLevelEmployees);

            callback(trees, employeesWithInvalidManager);
        });


    },

    printTree: function (t) {
        var output = [];
        function walk(node, currentDepth, maxDepth) {
            var row = [];
            for (let index = 0; index < maxDepth; index++) {
                if (currentDepth === index) {
                    row.push(node.value.name);
                } else {
                    row.push('');
                }
            }
            output.push(row);
            node.children.forEach(c => walk(c, currentDepth + 1, maxDepth));
            return output;
        }

        // Print the tree as a table
        var d = this.getDepth(t);
        walk(t, 0, d);
        return output;
    }
}
