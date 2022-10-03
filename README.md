# steampipe-plugin-ndo


![image](https://hub.steampipe.io/images/plugins/turbot/aws-social-graphic.png)

# NDO Plugin for Steampipe

Use SQL to query NDO infrastructure

- **[Get started →](https://hub.steampipe.io/plugins/justlikeef/ndo)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/justlikeef/ndo/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/justlikeef/steampipe-plugin-ndo/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install ndo
```

Run a query:

```sql
select * from aci_access_pol_leaf_int_pg_vpc_bundle_grp
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/justlikeef/steampipe-plugin-aci.git
cd steampipe-plugin-aci
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/ndo.spc
```

Try it!

```
steampipe query
> .inspect ndo
```

Further reading:

- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-aws/blob/main/LICENSE).

`help wanted` issues:

- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [AWS Plugin](https://github.com/turbot/steampipe-plugin-aws/labels/help%20wanted)
