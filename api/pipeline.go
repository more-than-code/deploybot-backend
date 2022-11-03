package api

import (
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/more-than-code/deploybot/model"
	"github.com/more-than-code/deploybot/repository"
)

type Api struct {
	repo *repository.Repository
}

type TaskCollection struct {
	Pipelines []*model.Pipeline
}

func NewApi() *Api {
	r, err := repository.NewRepository()
	if err != nil {
		panic(err)
	}
	return &Api{repo: r}
}

func (a *Api) DashboardHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pls, _ := a.repo.GetPipelines(ctx)

		tmpl := template.Must(template.ParseFiles("asset/tasks.html"))

		tmpl.Execute(ctx.Writer, pls)
	}
}

func (a *Api) PostPipeline() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.CreatePipelineInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostPipelineResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		id, err := a.repo.CreatePipeline(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PostPipelineResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PostPipelineResponse{Payload: PostPipelineResponsePayload{id}})
	}

}

func (a *Api) GetPipelines() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pls, err := a.repo.GetPipelines(ctx)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetPipelinesResponse{Code: CodeClientError, Msg: err.Error()})
		}

		ctx.JSON(http.StatusOK, GetPipelinesResponse{Payload: GetPipelinesResponsePayload{pls}})
	}
}

func (a *Api) GetPipeline() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Param("name")
		input := model.GetPipelineInput{Name: name}
		pl, err := a.repo.GetPipeline(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, GetPipelineResponse{Code: CodeClientError, Msg: err.Error()})
		}

		ctx.JSON(http.StatusOK, GetPipelineResponse{Payload: GetPipelineResponsePayload{pl}})
	}
}

func (a *Api) PatchPipeline() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.UpdatePipelineInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchPipelineResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdatePipeline(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PatchPipelineResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PatchPipelineResponse{})
	}
}

func (a *Api) PutPipelineStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var input model.UpdatePipelineStatusInput
		err := ctx.BindJSON(&input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PutPipelineStatusResponse{Code: CodeClientError, Msg: err.Error()})
			return
		}

		err = a.repo.UpdatePipelineStatus(ctx, &input)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, PutPipelineStatusResponse{Code: CodeServerError, Msg: err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, PutPipelineStatusResponse{})
	}
}