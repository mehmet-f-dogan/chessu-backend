package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"mehmetfd.dev/chessu-backend/database"
	"mehmetfd.dev/chessu-backend/models"

	"github.com/google/uuid"
)

func AssignHomepageHandlers(app *fiber.App) {
	app.Get("/homepage/user/:userId/courses", handleUserCourses)
}

type UserHomepageCoursesResponseItem struct {
	CourseId             string `json:"course_id"`
	CompletionPercentage uint8  `json:"completion"`
}

type UserHomepageCoursesResponse struct {
	CourseResponses []UserHomepageCoursesResponseItem
}

func handleUserCourses(c *fiber.Ctx) error {
	clerkUserId := utils.CopyString(c.Params("userId"))

	var user models.AppUser
	if err := database.DB.Where(&models.AppUser{ClerkId: clerkUserId}).First(&user).Error; err != nil {
		return c.JSON([]bool{})
	}

	// Mark the content as completed
	ids := [][16]byte{}
	for _, element := range user.PurchasedCourseId.Elements {
		ids = append(ids, element.Bytes)
	}
	for _, element := range user.CompletedContentId.Elements {
		coursePtr, _, _ := database.GetCourseAndChapterAndContent(element.Bytes)
		ids = append(ids, coursePtr.Id.Bytes)
	}

	ids = removeDuplicates(ids)

	responses := make([]UserHomepageCoursesResponseItem, len(ids))

	for i, courseId := range ids {
		courseIdStr, _ := uuid.FromBytes(courseId[:])
		responses[i] = UserHomepageCoursesResponseItem{
			CourseId:             courseIdStr.String(),
			CompletionPercentage: calculateCompletionPercentage(courseId, user.Id.Bytes),
		}
	}

	var response = UserHomepageCoursesResponse{
		CourseResponses: responses,
	}
	return c.JSON(response)
}

func removeDuplicates(arr [][16]byte) [][16]byte {
	// Create a map to store unique elements
	unique := make(map[[16]byte]bool)

	// Create a new slice to store unique elements
	result := [][16]byte{}

	// Iterate over the input array
	for _, num := range arr {
		// Check if the element is already in the map
		if !unique[num] {
			// Add the element to the map and slice
			unique[num] = true
			result = append(result, num)
		}
	}

	return result
}

func calculateCompletionPercentage(courseId uuid.UUID, userId uuid.UUID) uint8 {
	coursePtr := database.GetCourse(courseId)
	if coursePtr == nil {
		return 0
	}

	var user models.AppUser

	if err := database.DB.Find(&user, userId).Error; err != nil {
		return 0
	}

	completedContentIds := make([]uuid.UUID, 0)
	notCompletedContentIds := make([]uuid.UUID, 0)

	for _, chapter := range coursePtr.Chapters {
		for _, content := range chapter.Contents {
			found := false
			for _, completedContentId := range user.CompletedContentId.Elements {
				if completedContentId == content.Id.UUID {
					found = true
					break
				}
			}
			if found {
				completedContentIds = append(completedContentIds, content.Id.Bytes)
			} else {
				notCompletedContentIds = append(notCompletedContentIds, content.Id.Bytes)
			}
		}
	}
	if len(completedContentIds) == 0 && len(notCompletedContentIds) == 0 {
		return 100
	}
	percentage := uint8(float32(len(completedContentIds)) * 100 / float32(len(completedContentIds)+len(notCompletedContentIds)))
	return percentage
}
