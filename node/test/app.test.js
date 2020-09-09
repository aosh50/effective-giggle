const assert = require('assert');
const expect = require('chai').expect;
var app = require('../app.js');

describe('CSV File Load', async () => {
    it('should load 6 employees', () => {
        app.loadEmployeesFromFile("./data/employee_list.csv", employees => {
            assert.equal(employees.length, 6);
        });
    });

    it('should load the employees with the correct format', () => {
        const expected = [  new app.Employee('100', 'Alan', '150'),
            new app.Employee('220', 'Martin', '100')];

        app.loadEmployeesFromFile("./data/employee_list_short.csv", employees => {
            
            expect(employees).to.deep.equal(expected);
        });
    });

});

describe('Build Hierarchy', async () => {
    it('builds the hierarchy for the example problem', () => {
        const expected = [ 
            [ 'Jamie', '', '' ],
            [ '', 'Alan', '' ],
            [ '', '', 'Martin' ],
            [ '', '', 'Alex' ],
            [ '', 'Steve', '' ],
            [ '', '', 'David' ] 
        ];
        app.buildHierarchicalData("./data/employee_list.csv", trees => {
            trees.forEach(tree => {
                var result = app.printTree(tree);
                expect(result).to.deep.equal(expected)
            });

        }, err => {});
    });
    it('builds the hierarchy with extra levels', () => {
        const expected = [ 
            [ 'Jamie', '', '', '', '' ],
            [ '', 'Alan', '', '', '' ],
            [ '', '', 'Martin', '', '' ],
            [ '', '', 'Alex', '', '' ],
            [ '', 'Steve', '', '', '' ],
            [ '', '', 'David', '', '' ],
            [ '', '', '', 'John', '' ],
            [ '', '', '', '', 'Bob' ] 
        ];
        app.buildHierarchicalData("./data/employee_list_extra_levels.csv", trees => {
            trees.forEach(tree => {
                var result = app.printTree(tree);
                expect(result).to.deep.equal(expected)
            });

        }, err => {});
    });
    it('builds the hierarchy with multiple CEOs', () => {
        const expected = [ [ [ 'Jamie', '', '' ],
        [ '', 'Alan', '' ],
        [ '', '', 'Martin' ],
        [ '', '', 'Alex' ],
        [ '', 'Steve', '' ] ],
      [ [ 'David' ] ] ];
        app.buildHierarchicalData("./data/employee_list_two_ceos.csv", trees => {
            var results = [];
            trees.forEach(tree => {
                var result = app.printTree(tree);
                results.push(result);
            });
            expect(results).to.deep.equal(expected);

        }, err => {});
    });
    it('returns an error if no employees found', () => {
        app.buildHierarchicalData("./data/employee_list_empty.csv", trees => {}, 
        err => {
            expect(err).to.equal("No employees found!");
        });
    });
    it('returns an error if no CEOs present', () => {
        const expected = [ 
        ];
        app.buildHierarchicalData("./data/employee_list_no_ceos.csv", 
            trees => {}, 
            (err) => {
                expect(err).to.equal("No CEO found");
            }
        );
    });

    it('returns the employees with invalid managers', () => {
        const expected = [  new app.Employee('275', 'Alex', '103') ]
        app.buildHierarchicalData("./data/employee_list_invalid_manager.csv", (_, invalid) => {
            expect(invalid).to.deep.equal(expected)

        });
    });

});