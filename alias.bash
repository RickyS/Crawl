alias x='ls -lat *.txt'
alias g='grep -c "^Going" 1.txt | sed -e "s,^,Going ,"'
alias f='grep -c "^=Fetching=" 1.txt| sed -e "s,^,Fetch ,"'
alias e='grep -c "^REJECT" 1.txt | sed -e "s,^,Rejects ,"'
alias so='. ./alias.bash'
alias run='/bin/rm -f 1.txt 2.txt && go run crawl.go 1>1.txt 2>2.txt &'
alias k='kill %go'
alias 2='cat 2.txt'
alias 1='less 1.txt'

function creep () {
  egrep 'Creeping:' 1.txt
  egrep 'Current Domain: ' 2.txt
  x
  g
  f
  e
  grep -w NIL *.txt
  grep -c '^=ENDING' 1.txt | sed -e 's,^,Ending ,'
  grep -c 'Too many urls fetched' 1.txt | sed -e 's,^,Too many urls: ,'
  grep -c 'Test got result' 1.txt       | sed -e 's,^,Test got result: ,'
  egrep '^Waited:|Test Done|Test Closed| go Status: |StatusCode ' 1.txt
  grep -n DONE *.txt
  grep -i -w error *.txt
}

function scan ()  {
  #There's a better way to do this.
   grep -h routineStatus *.go | grep -v 'strings.Join' | sed -r -e 's,^[ \t]+,,' | sort
}