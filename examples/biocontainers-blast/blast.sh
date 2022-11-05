#!/bin/sh

# See: https://biocontainers-edu.readthedocs.io/en/latest/running_example.html
# cd /Users/yperez/workplace   # Replace by your path
# mkdir host-data
# docker run biocontainers/blast:2.2.31 blastp -help
# docker run -v `pwd`/host-data/ biocontainers/blast:2.2.31 curl -O ftp://ftp.ncbi.nih.gov/refseq/D_rerio/mRNA_Prot/zebrafish.1.protein.faa.gz
# docker run -v `pwd`/host-data/:/data/ biocontainers/blast:2.2.31 gunzip zebrafish.1.protein.faa.gz
# docker run -v `pwd`/host-data/:/data/ biocontainers/blast:2.2.31 makeblastdb -in zebrafish.1.protein.faa -dbtype prot
# docker run biocontainers/blast:2.2.31 curl https://www.uniprot.org/uniprot/P04156.fasta >> host-data/P04156.fasta
# docker run -v `pwd`/host-data/:/data/ biocontainers/blast:2.2.31 blastp -query P04156.fasta -db zebrafish.1.protein.faa -out results.txt

# we run all at once; just a shell script embedded in the 
# binary running in the container

blastp -help
curl -O ftp://ftp.ncbi.nih.gov/refseq/D_rerio/mRNA_Prot/zebrafish.1.protein.faa.gz
gunzip zebrafish.1.protein.faa.gz
makeblastdb -in zebrafish.1.protein.faa -dbtype prot
curl https://rest.uniprot.org/uniprotkb/P04156.fasta > P04156.fasta
echo "P04156.fasta"
cat P04156.fasta
echo "Running blastp"
blastp -query ./P04156.fasta -db zebrafish.1.protein.faa -out results.txt
echo "Results"
cat results.txt
echo "Copy results to host shared directory"
cp results.txt /host