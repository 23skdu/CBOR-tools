use Test::More;
$ENV{PATH} = "$ENV{PATH}:/bin:/usr/bin";
print "PreFlight Checks for Tests\n";
use_ok('Perl::Critic');
my $critic = Perl::Critic->new(-severity => 'stern');
my @flames=$critic->critique("$0");
print @flames;
ok ( (! @flames),'Perl::Critic has nothing to say');

