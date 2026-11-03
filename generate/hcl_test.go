package generate_test

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cycloidio/inframap/errcode"
	"github.com/cycloidio/inframap/generate"
	"github.com/cycloidio/inframap/graph"
)

func TestFromHCL_AWS(t *testing.T) {
	t.Run("SuccessSG", func(t *testing.T) {
		fs := afero.NewOsFs()

		g, cfg, err := generate.FromHCL(fs, "./testdata/aws_hcl_sg.tf", generate.Options{Clean: true, Connections: true, ExternalNodes: true})
		require.NoError(t, err)
		require.NotNil(t, g)

		eg := &graph.Graph{
			Nodes: []*graph.Node{
				&graph.Node{
					Canonical: "im_out.tcp/443->443",
				},
				&graph.Node{
					Canonical: "im_out.tcp/80->80",
				},
				&graph.Node{
					Canonical: "aws_lb.front",
					Name:      "some name",
				},
				&graph.Node{
					Canonical: "aws_launch_template.front",
				},
				&graph.Node{
					Canonical: "aws_db_instance.application",
				},
				&graph.Node{
					Canonical: "aws_elasticache_cluster.redis",
				},
			},
			Edges: []*graph.Edge{
				&graph.Edge{
					Source:     "im_out.tcp/443->443",
					Target:     "aws_lb.front",
					Canonicals: []string(nil),
				},
				&graph.Edge{
					Source:     "im_out.tcp/80->80",
					Target:     "aws_lb.front",
					Canonicals: []string(nil),
				},
				&graph.Edge{
					Source:     "aws_lb.front",
					Target:     "aws_launch_template.front",
					Canonicals: []string{"aws_security_group.lb-front", "aws_security_group.front"},
				},
				&graph.Edge{
					Source:     "aws_launch_template.front",
					Target:     "aws_db_instance.application",
					Canonicals: []string{"aws_security_group.front", "aws_security_group.rds"},
				},
				&graph.Edge{
					Source:     "aws_launch_template.front",
					Target:     "aws_elasticache_cluster.redis",
					Canonicals: []string{"aws_security_group.redis", "aws_security_group.front"},
				},
			},
		}

		assertEqualGraph(t, eg, g, cfg)
	})
}

func TestFromHCL_FlexibleEngine(t *testing.T) {
	t.Run("SuccessSG", func(t *testing.T) {
		fs := afero.NewOsFs()

		g, cfg, err := generate.FromHCL(fs, "./testdata/flexibleengine_hcl.tf", generate.Options{Clean: true, Connections: true, ExternalNodes: true})
		require.NoError(t, err)
		require.NotNil(t, g)

		eg := &graph.Graph{
			Nodes: []*graph.Node{
				&graph.Node{
					Canonical: "flexibleengine_compute_instance_v2.instance_one",
					Name:      "instance-one",
				},
				&graph.Node{
					Canonical: "flexibleengine_compute_instance_v2.instance_two",
					Name:      "instance-two",
				},
			},
			Edges: []*graph.Edge{
				&graph.Edge{
					Target: "flexibleengine_compute_instance_v2.instance_one",
					Source: "flexibleengine_compute_instance_v2.instance_two",
					Canonicals: []string{
						"flexibleengine_networking_port_v2.port_instance_two",
						"flexibleengine_networking_port_v2.port_instance_one",
						"flexibleengine_networking_secgroup_rule_v2.ssh_two_to_one",
						"flexibleengine_networking_secgroup_v2.secgroup_instance_two",
						"flexibleengine_networking_secgroup_v2.secgroup_instance_one",
					},
				},
			},
		}

		assertEqualGraph(t, eg, g, cfg)
	})
}

func TestFromHCL_Google(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		fs := afero.NewOsFs()

		g, cfg, err := generate.FromHCL(fs, "./testdata/google_hcl.tf", generate.Options{Clean: true, Connections: true, ExternalNodes: true})
		require.NoError(t, err)
		require.NotNil(t, g)

		eg := &graph.Graph{
			Nodes: []*graph.Node{
				&graph.Node{
					Canonical: "google_compute_instance.inframap-tmp-two",
					Name:      "inframap-tmp-two",
				},
				&graph.Node{
					Canonical: "google_compute_instance.inframap-tmp",
					Name:      "inframap-tmp",
				},
			},
			Edges: []*graph.Edge{
				&graph.Edge{
					Target: "google_compute_instance.inframap-tmp",
					Source: "google_compute_instance.inframap-tmp-two",
					Canonicals: []string{
						"google_compute_firewall.allow-ssh",
					},
				},
			},
		}

		assertEqualGraph(t, eg, g, cfg)
	})
}

func TestFromHCL_Module(t *testing.T) {
	t.Run("SuccessSG", func(t *testing.T) {
		fs := afero.NewOsFs()

		g, cfg, err := generate.FromHCL(fs, "./testdata/tf-module/", generate.Options{Clean: true, Connections: true, ExternalNodes: true})
		require.NoError(t, err)
		require.NotNil(t, g)

		eg := &graph.Graph{
			Nodes: []*graph.Node{
				&graph.Node{
					Canonical: "im_out.tcp/443->443",
				},
				&graph.Node{
					Canonical: "im_out.tcp/80->80",
				},
				&graph.Node{
					Canonical: "aws_lb.front",
					Name:      "some name",
				},
				&graph.Node{
					Canonical: "aws_launch_template.front",
				},
				&graph.Node{
					Canonical: "aws_db_instance.application",
				},
				&graph.Node{
					Canonical: "aws_elasticache_cluster.redis",
				},
			},
			Edges: []*graph.Edge{
				&graph.Edge{
					Source:     "im_out.tcp/443->443",
					Target:     "aws_lb.front",
					Canonicals: []string(nil),
				},
				&graph.Edge{
					Source:     "im_out.tcp/80->80",
					Target:     "aws_lb.front",
					Canonicals: []string(nil),
				},
				&graph.Edge{
					Source:     "aws_lb.front",
					Target:     "aws_launch_template.front",
					Canonicals: []string{"aws_security_group.lb-front", "aws_security_group.front"},
				},
				&graph.Edge{
					Source:     "aws_launch_template.front",
					Target:     "aws_db_instance.application",
					Canonicals: []string{"aws_security_group.front", "aws_security_group.rds"},
				},
				&graph.Edge{
					Source:     "aws_launch_template.front",
					Target:     "aws_elasticache_cluster.redis",
					Canonicals: []string{"aws_security_group.redis", "aws_security_group.front"},
				},
			},
		}

		assertEqualGraph(t, eg, g, cfg)
	})
}

func TestFromHCL_ModernLanguageFeatures(t *testing.T) {
	fs := afero.NewOsFs()

	g, _, err := generate.FromHCL(fs, "./testdata/tf-modern/", generate.Options{Raw: true})
	require.NoError(t, err)
	require.NotNil(t, g)

	bucket, err := g.GetNodeByCanonical("aws_s3_bucket.logs")
	require.NoError(t, err)
	policy, err := g.GetNodeByCanonical("aws_s3_bucket_policy.logs")
	require.NoError(t, err)

	edges := g.GetEdgesForNode(policy.ID)
	require.Len(t, edges, 1)
	assert.ElementsMatch(t, []string{bucket.ID, policy.ID}, []string{edges[0].Source, edges[0].Target})
}

func TestFromHCL_LoaderEdgeCases(t *testing.T) {
	const (
		bucket = "resource \"aws_s3_bucket\" \"a\" {\n  bucket = \"x\"\n}\n"
		policy = "resource \"aws_s3_bucket_policy\" \"a\" {\n  bucket = aws_s3_bucket.a.id\n}\n"
	)
	newFs := func(t *testing.T, files map[string]string) afero.Fs {
		fs := afero.NewMemMapFs()
		for name, src := range files {
			require.NoError(t, afero.WriteFile(fs, "mod/"+name, []byte(src), 0o644))
		}
		return fs
	}

	t.Run("OverrideFilesKeepBaseLinks", func(t *testing.T) {
		fs := newFs(t, map[string]string{
			"main.tf":          bucket + policy,
			"main_override.tf": "resource \"aws_s3_bucket_policy\" \"a\" {\n  policy = \"{}\"\n}\n",
			"override.tf":      "resource \"aws_s3_bucket\" \"a\" {\n  bucket = \"y\"\n}\n",
		})

		g, _, err := generate.FromHCL(fs, "mod", generate.Options{Raw: true})
		require.NoError(t, err)

		p, err := g.GetNodeByCanonical("aws_s3_bucket_policy.a")
		require.NoError(t, err)
		assert.Len(t, g.GetEdgesForNode(p.ID), 1)
	})

	t.Run("DuplicateResourceAcrossFiles", func(t *testing.T) {
		fs := newFs(t, map[string]string{"a.tf": bucket, "b.tf": bucket})

		_, _, err := generate.FromHCL(fs, "mod", generate.Options{Raw: true})
		require.Error(t, err)
		assert.Contains(t, err.Error(), `duplicate resource "aws_s3_bucket.a"`)
	})

	t.Run("JSONConfigurationIsRejected", func(t *testing.T) {
		fs := newFs(t, map[string]string{
			"main.tf":       policy,
			"extra.tf.json": `{"resource":{"aws_s3_bucket":{"a":{"bucket":"x"}}}}`,
		})

		_, _, err := generate.FromHCL(fs, "mod", generate.Options{Raw: true})
		assert.ErrorIs(t, err, errcode.ErrGenerateFromJSON)

		_, _, err = generate.FromHCL(fs, "mod/extra.tf.json", generate.Options{Raw: true})
		assert.ErrorIs(t, err, errcode.ErrGenerateFromJSON)
	})
}
