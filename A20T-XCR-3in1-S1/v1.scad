difference(){
    translate([-1,0,0]) cube([32,5,4.6]);

    translate([2.5,5.1,2.3]) rotate([90,0,0]) cylinder(h=7,r=1.5,$fn=100);
    translate([30-2.5,5.1,2.3]) rotate([90,0,0]) cylinder(h=7,r=1.5,$fn=100);

    translate([30-((30-18)/2),2.5,-1]) cylinder(h=7,r=1.5,$fn=100);
    translate([((30-18)/2),2.5,-1]) cylinder(h=7,r=1.5,$fn=100);
}