package provider

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/waduhek/flagger/internal/project"
)

type providerRepository struct {
	projectColl *mongo.Collection
}

func (r *providerRepository) GetFlagDetailsByProjectKey(
	ctx context.Context,
	projectKey string,
	environmentName string,
	flagName string,
) ([]FlagDetails, error) {
	pipeline := bson.A{
		bson.D{
			{
				Key:   "$match",
				Value: bson.D{{Key: "key", Value: projectKey}},
			},
		},
		bson.D{
			{
				Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "environments"},
					{Key: "localField", Value: "environments"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "environment"},
				},
			},
		},
		bson.D{
			{
				Key:   "$unwind",
				Value: bson.D{{Key: "path", Value: "$environment"}},
			},
		},
		bson.D{
			{
				Key: "$match",
				Value: bson.D{
					{Key: "environment.name", Value: environmentName},
				},
			},
		},
		bson.D{
			{
				Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "flags"},
					{Key: "localField", Value: "flags"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "flag"},
				},
			},
		},
		bson.D{
			{
				Key:   "$unwind",
				Value: bson.D{{Key: "path", Value: "$flag"}},
			},
		},
		bson.D{
			{
				Key:   "$match",
				Value: bson.D{{Key: "flag.name", Value: flagName}},
			},
		},
		bson.D{
			{
				Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "flag_settings"},
					{Key: "localField", Value: "flag_settings"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "flag_setting"},
				},
			},
		},
		bson.D{
			{
				Key:   "$unwind",
				Value: bson.D{{Key: "path", Value: "$flag_setting"}},
			},
		},
		bson.D{
			{
				Key: "$match",
				Value: bson.D{
					{
						Key: "$expr",
						Value: bson.D{
							{
								Key: "$and",
								Value: bson.A{
									bson.D{
										{
											Key: "$eq",
											Value: bson.A{
												"$flag_setting.project_id",
												"$_id",
											},
										},
									},
									bson.D{
										{
											Key: "$eq",
											Value: bson.A{
												"$flag_setting.environment_id",
												"$environment._id",
											},
										},
									},
									bson.D{
										{
											Key: "$eq",
											Value: bson.A{
												"$flag_setting.flag_id",
												"$flag._id",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{
				Key: "$project",
				Value: bson.D{
					{Key: "environments", Value: 0},
					{Key: "flags", Value: 0},
					{Key: "flag_settings", Value: 0},
				},
			},
		},
	}

	cursor, err := r.projectColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []FlagDetails
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func NewProviderRepository(db *mongo.Database) *providerRepository {
	projectColl := db.Collection(project.ProjectCollection)

	return &providerRepository{projectColl: projectColl}
}
