#!/usr/bin/env perl
use warnings; 
use strict;
use CBOR::Free; 
use Getopt::Std;
my %opt=();
getopts("f:a:t:s:o:",\%opt);
my $action=$opt{'a'};
my $file=$opt{'f'};
my $string=$opt{'s'};
my $output=$opt{'o'};
my $tag=$opt{'t'};
my @actions=qw/Encode Decode Tag/;
if (! defined $opt{'a'}){ print "ERR: Must specify action(-a)\n"; usage(); };
if (! grep(/^$action$/, @actions)) { print "ERR: Must specify a valid action\n"; usage(); }; 

eval { &{ \&$action }();};

sub usage { print "Actions: ";foreach my $entry(@actions) { print "$entry "};print "\n"; exit 1; }
sub Encode { }
sub Decode { }
sub Tag { }
