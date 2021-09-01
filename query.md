# ubuntu 
```
select cve_id,name,version,severity from packages,debians,advisories WHERE packages.definition_id = debians.definition_id and debians.definition_id = advisories.definition_id
```

# redhat
```
select cve_id,name,version,severity from packages,advisories,cves WHERE packages.definition_id = advisories.definition_id and advisories.id = cves.advisory_id

select cve_id,name,version from packages,advisories,cves WHERE packages.definition_id = advisories.definition_id and advisories.id = cves.advisory_id
```
# debian
```
select cve_id,name,version from packages,debians WHERE packages.definition_id = debians.definition_id
```