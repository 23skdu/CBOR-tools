#!/usr/bin/env perl
use warnings; 
use strict;
use LWP;
use Getopt::Std;
my %opt=();
getopts("a:",\%opt);
my $action=$opt{'a'};
my @actions=qw/SomeFunction/;
if (! defined $opt{'a'}){ print "ERR: Must specify action(-a)\n"; usage(); };
if (! grep(/^$action$/, @actions)) { print "ERR: Must specify a valid action\n"; usage(); }; 

eval { &{ \&$action }();};

sub usage { print "Actions: ";foreach my $entry(@actions) { print "$entry "};print "\n"; exit 1; }

sub SomeFunction { } 
