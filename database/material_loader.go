package database

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"

	"mehmetfd.dev/chessu-backend/models"
)

var Materials []models.Course = []models.Course{}

func LoadMaterials() {
	ctx := context.Background()
	config, err := config.LoadDefaultConfig(ctx,
		//config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_KEY"), os.Getenv("AWS_SECRET"), "")),
	)
	if err != nil {
		panic(err)
	}
	client := s3.NewFromConfig(config)
	bucketName := os.Getenv("AWS_MATERIALS_S3_BUCKET_NAME")

	listMaterialsParams := &s3.ListObjectsV2Input{
		Bucket: &bucketName,
	}

	material, err := client.ListObjectsV2(ctx, listMaterialsParams)
	if err != nil {
		panic(err)
	}

	for _, content := range material.Contents {
		contentData, err := client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucketName,
			Key:    content.Key,
		})
		if err != nil {
			panic(err)
		}
		contentBytes, err := io.ReadAll(contentData.Body)
		if err != nil {
			panic(err)
		}
		var c models.Course
		err = json.Unmarshal(contentBytes, &c)
		if err != nil {
			panic(err)
		}
		Materials = append(Materials, c)
	}
}

func GetCourse(courseId uuid.UUID) *models.Course {
	for _, course := range Materials {
		if course.Id.UUID.Bytes == courseId {
			return &course
		}
	}
	return nil
}

func GetCourseAndChapter(chapterId uuid.UUID) (*models.Course, *models.Chapter) {
	for _, course := range Materials {
		for _, chapter := range course.Chapters {
			if chapter.Id.UUID.Bytes == chapterId {
				return &course, &chapter
			}
		}
	}
	return nil, nil
}

func GetCourseAndChapterAndContent(contentId uuid.UUID) (*models.Course, *models.Chapter, *models.Content) {
	for _, course := range Materials {
		for _, chapter := range course.Chapters {
			for _, content := range chapter.Contents {
				if content.Id.UUID.Bytes == contentId {
					return &course, &chapter, &content
				}
			}
		}
	}
	return nil, nil, nil
}
